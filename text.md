1. Who uses this?

    The primary user is a sim racer that installs the ingest application/cli, it runs all the time and starts on windows boot automatically. After there race completes the .ibt file is created and the ingest will detect that and process the file.

    On the other end of that is the UI dashboard that allows sim racers as drivers and their team or performance engineer if they have one to analyse the racing/lap data comparing lap by lap data.

    Lets say there are 50 simultaneous users ingesting data and 70 simultaneous users using the ui at any time. The user will install some software on their sim racing pc which will be windows and will have a lot of power but we don't want to impact their sim racing fps.

  2. When is it used?

    For ease of use it waits for the sessions to be done and ibt file is made

  - Long after? (comparing this week vs last month, trend analysis)

  3. What decisions does it help make?

  - Examples: "brake 5m later into turn 3", "tire pressures are too high", "fuel won't last", "I'm faster in sector 2 than yesterday"
  - What question is the user asking when they open this tool?

    What speed did I carry at this apex. When did I get on throttle here in comparison to my team mate or another driver that saved data or my own other lap for now. Which sector was I faster or slower in

  4. What exists today?

  - What's working right now? (even if rough)

    The whole flow works right now but I want to push the ingest to see what is the fastest and best way to ingest data to the system using golang and the ibt files.

  - What do you currently use instead? (iRacing's built-in, motec, VRS, pen and paper, nothing?)

  VRS or garage61

  - What's frustrating about those alternatives?

    Nothing I just want to build my own to learn, This won't be a product for actual users

  5. Where does it run?

  - The iRacing machine: Do you control it? Is it always the same PC?

    Ingest runs on the sim racing machine. everything dockerised runs on an RPI 5 16gb for now but could run on a ryzen ai 7 370 if needed

  - The display: Same machine? Phone/tablet on a rig? Separate monitor? Remote?

    Sim rig so probalby a wide screen tripples or on a laptop

  - Connectivity: Local network only? Over the internet?

    tailscale/local for now

  6. Data sources

  - Live shared memory from iRacing — is this accessible to you?

    No

  - .ibt files only — are these being produced by your machine or uploaded by others?

    It's being made by iracing on the sim machine so I have old ones created from my sim rig but users would be submitting theirs created on there devices 

  - Other data? (video, setup sheets, standings from iRacing API)

  not yet but I could link to the iracing standings page for each session 

  7. Non-negotiables vs nice-to-haves

    Not lose or degrade the data
    Be fast and effient
    The UI must be comprehensive to allow race engineers to break down the data into usable outcomes
  List 2-3 things it must do on day one, and anything else you've been imagining but could live without initially

  8. What is the object as the builder of this project

    To replicate real motorsports systems to gain a better insight into that world whilst using the best of my skills, tools and knowledge. I have some data from sim racing which is a good starting point but I want to take that and use it to build more systems and gain insight into what uses formula one, world endurance championship, imsa, nascar, dtm and more have for software and software engineers