# SE Project : Sprint 4

This document contains the details of the work done in Sprint 4 of the Software Engineering Project. This sprint focused on implementing cross-platform window observation capabilities, improving documentation, and enhancing the overall system architecture. Building upon Sprint 3's foundation of testing and core functionality, this sprint introduced significant platform-specific features and improved system documentation.

Repository Link: [Timelygator/Timelygator](https://github.com/timelygator/TimelyGator)
Sprint4.md: [docs/files/sprint/sprint4.md](https://github.com/timelygator/TimelyGator/blob/main/docs/files/sprint/Sprint4.md)

## Sprint 3 Rolled Over Tasks

The following tasks from Sprint 3 were completed in Sprint 4:

### Frontend

- Continued improvements to UsageOverviewChart component
- Fixed UI/UX inconsistencies and theme persistence issues

### Backend

- Enhanced browser-observer implementation with cross-platform support
- Improved test coverage and documentation

## New Functionality

### Frontend

- Updated UsageOverviewChart component with improved data visualization
- Field modifications to extract specefic element from API calls
- Fixed UI/UX issues for better user experience
- Added FrontEnd documentation

### Backend

- Implemented cross-platform window observer for Linux, Windows, and macOS
- Added platform-specific window management using X11, Windows API, and macOS JXA
- Added AppleScript and JXA scripts for retrieving front application and window title
- Implemented proper permissions handling for macOS accessibility access
- Added CORS headers to allow cross-network access
- Fixed duration handling to use float64 seconds consistently
- Updated AFK timeout and poll time defaults for improved configuration

## Documentation

### API Documentation

This can now be visited at [timelygator.github.io/TimelyGator](https://timelygator.github.io/TimelyGator/), including architecture details, sprint information, etc.

- Implemented MkDocs for GitHub Pages documentation
- Added Swagger.json integration with MkDocs
- Moved sprint documentation to the website
- Updated documentation index to reflect README details
- Added comprehensive FE and BE documentation

## Tasks

1. [@siddhant-0707](https://github.com/siddhant-0707) # Backend
    - Implemented cross-platform window observer
    - Added platform-specific window management
    - Fixed duration handling and configuration
    - Added macOS permissions handling for accessibility access
    - Wrote documentation for observers and client

2. [@shreyansh-nayak-ufl](https://github.com/shreyansh-nayak-ufl) # Backend
    - Implemented MkDocs for documentation
    - Fixed CORS and network access issues
    - Added Linux stub implementation
    - Wrote documentation for FE/BE Architecture
    - Deployed GitHub Actions for automated docs website

3. [@PulkitGarg777](https://github.com/PulkitGarg777) # Frontend
    - Continued improvements to UsageOverviewChart component
    - Fixed UI/UX inconsistencies and ensured theme persistence across sessions
    - Wrote unit tests for UsageOverviewChart and other key UI elements
    - Worked on FrontEnd documentation

4. [@YashDVerma](https://github.com/YashDVerma) # Frontend
    - Testing API call over UsageOverviewChart 
    - Added Profile editing feature
    - Added FrontEnd documentation
    - Fixed ActualCharts not reacting to changed TimeRange

## Conclusion

Sprint 4 successfully implemented cross-platform window observation capabilities and significantly improved the project's documentation. The team focused on platform-specific implementations, system architecture improvements, and comprehensive documentation. The application now supports window tracking across all major operating systems and has improved documentation for better maintainability and user understanding.
