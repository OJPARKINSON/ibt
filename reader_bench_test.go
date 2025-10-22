package ibt

import (
	"os"
	"testing"
)

// BenchmarkMmapReaderCreation measures the cost of opening and mapping files
func BenchmarkMmapReaderCreation(b *testing.B) {
	testFile := ".testing/valid_test_file.ibt"

	b.Run("NewMmapReader", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			reader, err := NewMmapReader(testFile)
			if err != nil {
				b.Fatalf("failed to create mmap reader: %v", err)
			}
			reader.Close()
		}
	})
}

// BenchmarkMmapReaderData measures the cost of accessing the entire mmap'd data
func BenchmarkMmapReaderData(b *testing.B) {
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer reader.Close()

	b.Run("Data_access", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			data := reader.Data()
			if len(data) == 0 {
				b.Fatal("empty data")
			}
		}
	})
}

// BenchmarkMmapReaderRead measures the Read() method performance
func BenchmarkMmapReaderRead(b *testing.B) {
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer reader.Close()

	buf := make([]byte, 4096)

	b.Run("Read_4KB", func(b *testing.B) {
		b.SetBytes(4096)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMmapReaderReadAt measures ReadAt() performance at various offsets
func BenchmarkMmapReaderReadAt(b *testing.B) {
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer reader.Close()

	buf := make([]byte, 1024)
	data := reader.Data()
	fileSize := int64(len(data))

	b.Run("ReadAt_beginning", func(b *testing.B) {
		b.SetBytes(1024)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
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
		for i := 0; i < b.N; i++ {
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
		for i := 0; i < b.N; i++ {
			offset := offsets[i%len(offsets)]
			_, err := reader.ReadAt(buf, offset)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMmapReaderReadAtUnsafe measures zero-copy unsafe reads
func BenchmarkMmapReaderReadAtUnsafe(b *testing.B) {
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer reader.Close()

	b.Run("ReadAtUnsafe_4KB", func(b *testing.B) {
		b.SetBytes(4096)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data := reader.ReadAtUnsafe(0, 4096)
			if len(data) != 4096 {
				b.Fatal("wrong size")
			}
		}
	})

	b.Run("ReadAtUnsafe_100KB", func(b *testing.B) {
		size := 100 * 1024
		b.SetBytes(int64(size))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data := reader.ReadAtUnsafe(0, size)
			if len(data) != size {
				b.Fatal("wrong size")
			}
		}
	})
}

// BenchmarkMmapReaderClose measures cleanup performance
func BenchmarkMmapReaderClose(b *testing.B) {
	b.Run("Close", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			reader, err := NewMmapReader(".testing/valid_test_file.ibt")
			if err != nil {
				b.Fatalf("failed to create reader: %v", err)
			}
			b.StartTimer()

			err = reader.Close()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMmapReaderEmptyFile tests edge case performance
func BenchmarkMmapReaderEmptyFile(b *testing.B) {
	// Create empty test file
	emptyFile := ".testing/empty_bench.ibt"
	f, err := os.Create(emptyFile)
	if err != nil {
		b.Fatal(err)
	}
	f.Close()
	defer os.Remove(emptyFile)

	b.Run("EmptyFile_creation", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			reader, err := NewMmapReader(emptyFile)
			if err != nil {
				b.Fatal(err)
			}
			reader.Close()
		}
	})
}

// BenchmarkMmapReaderVsOsRead compares mmap vs traditional file reading
func BenchmarkMmapReaderVsOsRead(b *testing.B) {
	testFile := ".testing/valid_test_file.ibt"

	b.Run("MmapReader_sequential_read", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			reader, err := NewMmapReader(testFile)
			if err != nil {
				b.Fatal(err)
			}
			data := reader.Data()
			// Simulate reading all data
			sum := 0
			for _, b := range data {
				sum += int(b)
			}
			reader.Close()
		}
	})

	b.Run("os.ReadFile_sequential_read", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
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
