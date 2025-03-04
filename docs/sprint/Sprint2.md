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
  - Implement the `Client` package to handle API requests.
- `AFK Observer`
- `Chrome Browser Observer`
- Testing

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
- Unit tests for `Routes` under the `api` package to validate the response format and error handling.

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
    - 

4. [@siddhant-0707](https://github.com/siddhant-0707) # Backend
    - 

## Conclusion

The tasks for Sprint 2 were completed successfully. The frontend and backend were integrated, and the necessary unit tests were written. The team will now focus on achieving full integration between the observers and the web UI in the next sprint. End-to-end (E2E) testing will be conducted to ensure seamless user workflows. Documentation for the API and key components will also be written.