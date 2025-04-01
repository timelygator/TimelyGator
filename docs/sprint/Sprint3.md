# SE Project : Sprint 3

This document contains the details of the work done in Sprint 3 of the Software Engineering Project. This sprint focused on completing unfinished tasks from Sprint 2, implementing new functionality, and comprehensive testing of all components. The main goal was to ensure the application is ready with proper test coverage and documentation.

Repository Link: [Timelygator/Timelygator](https://github.com/timelygator/TimelyGator)
Sprint3.md: [docs/sprint/sprint3.md](https://github.com/timelygator/TimelyGator/blob/main/docs/sprint3.md)

Table of contents:

- [Sprint 2 Rolled Over Tasks](#sprint-2-rolled-over-tasks)
- [New Functionality](#new-functionality)
- [Testing](#testing)
  - [Frontend Tests](#frontend-tests)
  - [Backend Tests](#backend-tests)
- [Documentation](#documentation)
- [Tasks](#tasks)
- [Conclusion](#conclusion)

## Sprint 2 Rolled Over Tasks

The following tasks from Sprint 2 were completed in Sprint 3:

### Frontend
- Component testing left in sprint 2
- File cleanups for beter code quality 

### Backend
- Added `browser-observer` tree to `server/observers` code. [PR #98](https://github.com/timelygator/TimelyGator/pull/98)
- Moved all API test cases to centralized file under `routes_test.go` [PR #100](https://github.com/timelygator/TimelyGator/pull/100)

## New Functionality

### Frontend
- Started off with a bit of e2e testing and some real-time updation feature retention
- Testing API call behaviours manually to identify trend of data fetching and amking modifcations accordingly
- Persistent theme changes allows user to change website aesthetics.

### Backend
- Utilize `chrome.storage.sync` for persistance in configuration of browser extension.

## Testing

### Frontend Tests

#### Component Tests
1. **Real time updation testing**
   - Test real-time data updates

2. **API behaviour analysis**
   - Test WebSocket connections and API fetch behaviour

3. **UI Improvements**
   - Test responsive design and other minor UI improvements 

#### End-to-End Tests
- Complete user registration flow ,some tasks will be completed in later sprint but main goal is to produce good API fetches

### Backend Tests

#### `tg-fakedata` Tests
1. **`TestParseDateFlag`**  
   - Verifies that `parseDateFlag()` correctly parses valid date strings into `time.Time` objects.  
   - Ensures that invalid formats return appropriate errors.

2. **`TestSameDay`**  
   - Tests the `sameDay()` function to confirm it accurately detects if two timestamps fall on the same calendar day in UTC.  
   - Validates both positive (same day) and negative (different days) scenarios.

3. **`TestWeightedChoice`**  
   - Ensures `weightedChoice()` returns items based on defined probability weights.  
   - Performs statistical assertions over multiple iterations to check that higher-weighted items appear more frequently.

4. **`TestPickDuration`**  
   - Validates that `pickDuration()` returns a duration within an expected range when a base minute value is provided.  
   - Tests fallback behavior with `minutes == 0`, confirming it generates durations between 5 seconds and the max allowed seconds.

5. **`TestGetString`**  
   - Confirms that `getString()` extracts string values correctly from a JSON-encoded data blob.  
   - Ensures it returns an empty string when a key is not present in the JSON.

**Chrome Extension Tests**
   - Test Chrome extension integration

## Documentation

### API Documentation
- Updated Swagger documentation with new endpoints for `heartbeat`.

## Tasks

1. [@PulkitGarg777](https://github.com/PulkitGarg777) # Frontend
    - Implemented real-time data updates
    - Setup basic e2e testing files for some components
    - API fetch behaviour testing for Tabs Table


2. [@YashDVerma](https://github.com/YashDVerma) # Frontend
    - Fixed UI/UX inconsistencies
    - Fixed theme persistence issues
    - added unit testing for new features


3. [@shreyansh-nayak-ufl](https://github.com/shreyansh-nayak-ufl) # Backend
    - Completed Chrome extension implementation
    - Implemented `chrome.storage.sync` API
    - Unified BE test cases under `routes_test.go`


4. [@siddhant-0707](https://github.com/siddhant-0707) # Backend
    - Developed `tg-fakedata` to generate synthetic event data for TimelyGator
    - Included features for simulating AFK status, window focus, and browser activity
    - Wrote tests for `tg-fakedata`


## Conclusion

Sprint 3 successfully completed the remaining tasks from Sprint 2 and implemented new functionality to enhance the application. The team focused on comprehensive testing, documentation, and performance optimization. The application now has proper test coverage, improved user experience, and robust backend functionality.
