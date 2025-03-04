# SE Project : Sprint 2

This document contains the details of the work done in Sprint 2 of the Software Engineering Project. This sprint we focused on integrating the frontend with the backend, setting up the necessary dependencies, and writing unit tests for the frontend and backend components. The main goal was to ensure that the frontend and backend work together seamlessly.

Repository Link: [Timelygator/Timelygator](https://github.com/timelygator/TimelyGator)
Sprint2.md: [docs/sprint/sprint2.md](https://github.com/timelygator/TimelyGator/blob/main/docs/sprint2.md)

Table of contents:

- [SE Project : Sprint 2](#se-project--sprint-2)
  - [Sprint 1 Rolled Over Tasks](#sprint-1-rolled-over-tasks)
  - [User Stories](#user-stories)
  - [Unit Tests](#unit-tests)
    - [Frontend](#frontend)
    - [Backend](#backend)
  - [Tasks](#tasks)
  - [Conclusion](#conclusion)

## Sprint 1 Rolled Over Tasks

Some of the tasks from Sprint 1 were carried over to Sprint 2. These tasks were completed in Sprint 2.

- Frontend
  - Address any pending UI/UX improvements.
  - Refactor code where necessary for better maintainability.
- Backend
  - Completed pending endpoints for `/heartbeat` and `/data`.
  - Ensure proper error handling and response formats.

## User Stories

### Frontend

Frontend Work Summary as mentioned in [Discussion #49](https://github.com/timelygator/TimelyGator/discussions/49)

- Working with APIs & Data Display
  - Fetch and display data correctly in respective fields or StatBox components.
  - Ensure proper error handling and loading states.
- Setting Up Dependencies
  - Install and configure required testing libraries.
  - Set up necessary state management or API-handling tools.
- Frontend-Backend Integration
  - Connect the frontend with the backend API.
  - Ensure proper data flow between the two.
  - Debug and fix any API-related issues.
- Testing
  - Write a simple test using Cypress to validate key workflows.
  - Implement unit tests for frontend components, aiming for a 1:1 test-to-function ratio.
- Additional Suggested Tasks:
  - Ensure consistent coding styles across the frontend.
  - Add documentation for key components and API usage.

### Backend

Backend Work Summary inlcudes completion of the following tasks,

- `Client` package
The `TimelyGatorClient` serves as the primary interface for interacting with the TimelyGator server. The following tasks were completed as part of the client implementation:
  - **GET, POST, DELETE Requests:** Added utility methods (`get`, `post`, `deleteReq`) for interacting with the server API.
  - **Bucket Management:**
    - Create, delete, and manage buckets on the server.
    - Integrated functionality to handle both queued and immediate requests.
  - **Event Handling:**
    - Implemented methods to retrieve, insert, and delete events from specific buckets.
    - Added support for retrieving event counts and exporting/importing event data.
  - **Heartbeat Functionality:**
    - Developed `Heartbeat` method to manage activity data through server heartbeats.
    - Included pre-merge and post-merge logic to optimize event data flow.
  - **Query Execution:**
    - Added `Query` method to send custom queries to the server and retrieve results.

- `AFK Observer`
The `AFKWatcher` module was developed to monitor user inactivity and send AFK status to the server through the client. The work included:
  - **Platform Validation:** Ensured that the `AFKWatcher` only runs on supported platforms (Linux).
  - **Signal Handling:** Added graceful shutdown capabilities using system signals (`os.Interrupt`, `syscall.SIGTERM`).
  - **AFK Detection Loop:**
    - Implemented the `heartbeatLoop()` to regularly check AFK status.
    - Utilized `SecondsSinceLastInput()` to detect inactivity.
  - **Heartbeat Event Creation:**
    - Created `Event` objects with status (`afk` or `not-afk`) as JSON.
    - Sent accurate heartbeat data to the server using the `client.Heartbeat()` method.

- `Chrome Browser Observer`
The Chrome Browser Extension for TimelyGator was developed to capture detailed browsing activity, including active tabs, URLs, and time spent on websites. The extension integrates seamlessly with the TimelyGator client to provide accurate and categorized tracking of web activity.
  - **URL Tracking**
     - Automatically captures the active tab's URL.
     - Sends browsing data to the TimelyGator server at regular intervals.
  - **Tab Activity Monitoring**
     - Detects when a tab is opened, closed, or switched.
     - Tracks time spent on each tab and categorizes the data.
  - **Event Communication**
     - Utilizes message passing between content scripts and the background script.
     - Establishes a secure and efficient connection with the TimelyGator client
  - **Manifest v3 Compatibility:**
    - Built using the latest Chrome extension guidelines and technologies.
    - Ensures compliance with Chrome Web Store policies.

- Swagger API Documentation
  - **API Documentation:** Generated detailed documentation for the TimelyGator API using Swagger.
    - Used the `swag` package to automatically generate API documentation.
    - Ensured that all API endpoints, request/response formats, and error codes were documented.
  - **API Testing:** Validated the API documentation by testing each endpoint using Swagger UI.
    - Verified that the API responses matched the expected formats.
    - Ensured that the API was fully functional and ready for integration with the frontend.
  - **Postman Collection:** Created a Postman collection for testing the TimelyGator API.
    - Included sample requests for each endpoint to facilitate testing and debugging.
    - Ensured that the Postman collection was up-to-date with the latest API changes.

## Unit Tests

The following unit tests were written for the frontend and backend components.

### Frontend
- Configured Cypress testing and created related `.cy` files to achieve 1:1 testing, covering:
  1. Mounting of all components.
  2. Sidebar collapsing and expanding behavior.
  3. On-hover functionality of StatCard.
  4. Line graph, bar graph, and pie chart animations.
  5. Search functionality of TabTable.
  6. Toggle functionality on the settings page.
  7. Button click functionality on the settings page.

### Backend

- Unit tests for the `Client` package to ensure proper API request handling.

1. **`TestGetInfo`**
   - Verifies that `GetInfo()` correctly fetches server information and returns the expected status.
2. **`TestGetEvent`**
   - Tests the `GetEvent()` method to ensure the correct retrieval of a specific event by ID.
   - Simulates a server response with event data and checks if the event ID matches.
3. **`TestGetEvents`**
   - Confirms that `GetEvents()` fetches a list of events with the expected length and structure.
4. **`TestInsertEvent`**
   - Validates that `InsertEvent()` successfully posts a single event to the server without errors.
5. **`TestInsertEvents`**
   - Similar to `TestInsertEvent` but for multiple events, verifying batch event insertion.
6. **`TestDeleteEvent`**
   - Ensures `DeleteEvent()` properly handles event deletion requests and server responses.
7. **`TestGetEventCount`**
   - Checks that `GetEventCount()` returns the correct number of events as indicated by the mock server.
8. **`TestHeartbeat`**
   - Evaluates the `Heartbeat()` method in both queued and non-queued modes.
   - Simulates a "heartbeat received" response from the server.
9. **`TestGetBucketsMap`**
   - Tests the `GetBucketsMap()` method to ensure it correctly retrieves the list of buckets from the server.
10. **`TestCreateAndDeleteBucket`**
    - Combines testing for `CreateBucket()` and `DeleteBucket()` methods.
    - Simulates the full lifecycle of bucket creation and deletion.
11. **`TestExportAll`**
    - Validates that `ExportAll()` properly fetches all export data from the server.
12. **`TestExportBucket`**
    - Similar to `TestExportAll`, but specific to exporting data from a single bucket.
    - Confirms that the returned data matches the expected "bucket data" value.
13. **`TestImportBucket`**
    - Checks that `ImportBucket()` correctly posts import data to the server and handles the response.
14. **`TestQuery`**
    - Tests the `Query()` method by sending a sample query with time periods and validating the result count.
    - Simulates a "result1" and "result2" response from the server.
15. **`TestGetAndSetSetting`**
    - Covers both `GetSetting()` and `SetSetting()` methods.
    - Verifies that settings are correctly retrieved and updated.
16. **`TestWaitForStart`**
    - Simulates server readiness using `WaitForStart()`.
    - Tests both immediate success and timeout scenarios.

- Unit tests for `Routes` under the `api` package to validate the response format and error handling.

1. **`TestGetInfo`**
    – Verifies that `GetInfo()` correctly fetches server information (hostname, version, server_name) from the API.
2. **`TestGetBuckets`**
    – Checks that `GetBuckets()` returns a map with all expected buckets and proper bucket details.
3. **`TestGetBucket`**
    – Confirms that `CreateBucket()` accepts the correct payload and returns a successful response when creating a new bucket.
4. **`TestUpdateBucket`**
    – Validates that `UpdateBucket()` properly updates the bucket information and returns success status.
5. **`TestDeleteBucket`**
    – Ensures that `DeleteBucket()` with a force flag correctly deletes the specified bucket.
6. **`TestGetEvents`**
    – Verifies that `GetEvents()` returns a list of events with correct details for a given bucket.
7. **`TestGetEvent`**
    – Checks that `ExportAll()` returns all the export data, including all buckets with their associated details.
8. **`TestExportBucket`**
    – Confirms that `ImportAll()` successfully posts bucket data to the API and handles the server response.

## Tasks

TO get a better understanding of the user stories, their status, and people contributing to them,

1. [@PulkitGarg777](https://github.com/PulkitGarg777) # Frontend
    - Cleaned up some leftover codebases from Sprint1
    - Create API response format list needed for data fetching
    - Modify files to hold API templates for integration
    - Configure Cypress unit testing
    - Test 1:1 functionalities using Cypress
[issue #85](https://github.com/timelygator/TimelyGator/issues/85)
[issue #53](https://github.com/timelygator/TimelyGator/issues/53)
[PR #72](https://github.com/timelygator/TimelyGator/pull/72)
[PR #84](https://github.com/timelygator/TimelyGator/pull/84)

2. [@YashDVerma](https://github.com/YashDVerma) # Frontend
    -

3. [@shreyansh-nayak-ufl](https://github.com/shreyansh-nayak-ufl) # Backend
    - Worked on the `Browser Observer` module.
    - Created a Manifest v3 compatible Chrome extension.
    - Configured message passing between content scripts and background script.
    - Implemented URL tracking and tab activity monitoring.
    - Did complete testing of the `Server` routes.
    - Created a Postman collection for testing the TimelyGator API.
    - Generated detailed API documentation using Swagger.
[issue #69](https://github.com/timelygator/TimelyGator/issues/69)
[issue #70](https://github.com/timelygator/TimelyGator/issues/70)
[issue #71](https://github.com/timelygator/TimelyGator/issues/71)
[issue #19](https://github.com/timelygator/TimelyGator/issues/19)

4. [@siddhant-0707](https://github.com/siddhant-0707) # Backend
    - Implemented a `Client` module to interact with the TimelyGator `Server`
    - Built `AFK-Observer` which detects AFK status and sends heartbeats to the server.
    - Implemented keyboard and mouse listeners to capture user activity contributing to AFK detection.
    - Implemented API requests for data transfer between `AFK-Observer` and `Server` through the `Client`.
    - Wrote comprehensive unit tests for `Client`.
[issue #50](https://github.com/timelygator/TimelyGator/issues/50)
[issue #74](https://github.com/timelygator/TimelyGator/issues/74)
[PR #77](https://github.com/timelygator/TimelyGator/pull/77)
[PR #65](https://github.com/timelygator/TimelyGator/pull/65)

## Conclusion

The tasks for Sprint 2 were completed successfully. The frontend and backend were integrated, and the necessary unit tests were written. The team will now focus on achieving full integration between the observers and the web UI in the next sprint. End-to-end (E2E) testing will be conducted to ensure seamless user workflows. Documentation for the API and key components will also be written.
