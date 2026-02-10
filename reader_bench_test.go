package ibt

import (
	"os"
	"testing"
)

const benchTestFile = ".testing/valid_test_file.ibt"

// BenchmarkMmapReaderCreation measures the cost of opening and mapping files.
func BenchmarkMmapReaderCreation(b *testing.B) {
	testFile := benchTestFile

	b.Run("NewMmapReader", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			reader, err := NewIbtReader(testFile)
			if err != nil {
				b.Fatalf("failed to create mmap reader: %v", err)
			}
			_ = reader.Close()
		}
	})
}

// BenchmarkMmapReaderData measures the cost of accessing the entire mmap'd data.
func BenchmarkMmapReaderData(b *testing.B) {
	reader, err := NewIbtReader(benchTestFile)
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer func() { _ = reader.Close() }()

	b.Run("Data_access", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			data := reader.Data()
			if len(data) == 0 {
				b.Fatal("empty data")
			}
		}
	})
}

// BenchmarkMmapReaderRead measures the Read() method performance.
func BenchmarkMmapReaderRead(b *testing.B) {
	reader, err := NewIbtReader(benchTestFile)
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer func() { _ = reader.Close() }()

	buf := make([]byte, 4096)

	b.Run("Read_4KB", func(b *testing.B) {
		b.SetBytes(4096)
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_, err := reader.Read(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMmapReaderReadAt measures ReadAt() performance at various offsets.
func BenchmarkMmapReaderReadAt(b *testing.B) {
	reader, err := NewIbtReader(benchTestFile)
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer func() { _ = reader.Close() }()

	buf := make([]byte, 1024)
	data := reader.Data()
	fileSize := int64(len(data))

	b.Run("ReadAt_beginning", func(b *testing.B) {
		b.SetBytes(1024)
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_, err := reader.ReadAt(buf, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ReadAt_middle", func(b *testing.B) {
		offset := fileSize / 2
		b.SetBytes(1024)
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_, err := reader.ReadAt(buf, offset)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ReadAt_random_seeks", func(b *testing.B) {
		offsets := make([]int64, 100)
		for i := range offsets {
			offsets[i] = int64(i * 1000 % int(fileSize-1024))
		}
		b.SetBytes(1024)
		b.ReportAllocs()
		b.ResetTimer()
		for i := range b.N {
			offset := offsets[i%len(offsets)]
			_, err := reader.ReadAt(buf, offset)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkIbtReaderClose measures cleanup performance.
// Pre-creates readers in setup to avoid StopTimer/StartTimer hang.
func BenchmarkIbtReaderClose(b *testing.B) {
	const batchSize = 256
	readers := make([]*IbtReader, batchSize)
	for i := range readers {
		r, err := NewIbtReader(benchTestFile)
		if err != nil {
			b.Fatalf("failed to create reader: %v", err)
		}
		readers[i] = r
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		idx := i % batchSize
		if readers[idx] == nil {
			// Re-create if we've already closed this one
			r, err := NewIbtReader(benchTestFile)
			if err != nil {
				b.Fatalf("failed to create reader: %v", err)
			}
			readers[idx] = r
		}
		err := readers[idx].Close()
		if err != nil {
			b.Fatal(err)
		}
		readers[idx] = nil
	}
}

// BenchmarkMmapReaderEmptyFile tests edge case performance.
func BenchmarkMmapReaderEmptyFile(b *testing.B) {
	// Create empty test file
	emptyFile := ".testing/empty_bench.ibt"
	f, err := os.Create(emptyFile)
	if err != nil {
		b.Fatal(err)
	}
	_ = f.Close()
	defer func() { _ = os.Remove(emptyFile) }()

	b.Run("EmptyFile_creation", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			reader, err := NewIbtReader(emptyFile)
			if err != nil {
				b.Fatal(err)
			}
			_ = reader.Close()
		}
	})
}

// BenchmarkIbtReaderVsOsRead compares IbtReader vs raw os.ReadFile for full-file reads.
func BenchmarkIbtReaderVsOsRead(b *testing.B) {
	testFile := benchTestFile

	b.Run("IbtReader_sequential_read", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			reader, err := NewIbtReader(testFile)
			if err != nil {
				b.Fatal(err)
			}
			data := reader.Data()
			// Simulate reading all data
			sum := 0
			for _, b := range data {
				sum += int(b)
			}
			_ = reader.Close()
			_ = sum
		}
	})

	b.Run("os.ReadFile_sequential_read", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			data, err := os.ReadFile(testFile)
			if err != nil {
				b.Fatal(err)
			}
			// Simulate reading all data
			sum := 0
			for _, b := range data {
				sum += int(b)
			}
			_ = sum
		}
	})
}
