# SE Project : Sprint 1

This document contains the details of the first sprint of the Software Engineering project. The sprint focused on setting up the project structure, defining the main objectives, and evaluating feasibility. The goal of the sprint is to establish a solid foundation for the project and define the roadmap for future development.

Repository Link: [Timelygator/Timelygator](https://github.com/timelygator/TimelyGator)
Sprint1.md: [docs/sprint/sprint1.md](https://github.com/timelygator/TimelyGator/blob/main/docs/sprint1.md)

Table of contents:

- [Discussions](#discussion-1--project-structure)
  - [Project structure](#discussion-1--project-structure)
  - [Goals and competing products](#discussion-2--goals-and-competing-products)
  - [Tech stack](#discussion-3--tech-stack)
  - [Documentation](#discussion-4--documentation)
  - [Repository settings](#discussion-5--repository-settings)
- [User stories](#user-stories)
- [Tasks](#tasks)
- [Conclusion](#conclusion)

## Discussion 1 : Project structure

For the project structure, we decided following points,

- A monorepo approach with separate folders for frontend and backend codebases.
- Documentation will be hosted in a separate folder.
- Semantic versioning will be followed for versioning the project.
- GitHub discussions will be used for team communication and decision-making, and GitHub issues will be used for tracking tasks and bugs.

The project structure is as follows:

```project-name/
├── web-ui/              # Frontend codebase
├── server/              # Backend codebase
├── docs/                # Documentation
├── CONTRIBUTING.md      # Guidelines for contributing to the project.
├── CODE_OF_CONDUCT.md   # A code of conduct for team interactions.
├── LICENSE              # License file
├── CHANGELOG.md         # Document changes and updates to the project over time.
└── README.md            # Project README
```

> Links: [discussion #5](https://github.com/timelygator/TimelyGator/discussions/5) [issue #20](https://github.com/timelygator/TimelyGator/issues/20)

### Discussion 2 : Goals and competing products

The main objectives of the project are:

- To provide automated time tracking with minimal user input.
- Store time series data for analysis and reporting.
- Provide insights and recommendations based on user data.
- Support cross-platform usage with synchronization.
- Create a user-friendly interface for tracking time and setting goals.

We have also identified some competing products in the market:

- [Rize](https://www.rize.io/)
- [Toggl](https://toggl.com/)
- [ActivityWatch](https://activitywatch.net/)

### Discussion 3 : Tech stack

For the tech stack, we decided to use the following technologies:

- Frontend:
  - `react: 18.3.1`: JavaScript library for building user interfaces.
  - `react-dom: 18.3.1`: DOM-specific methods that can be used at the top level of a web app.
  - `react-router-dom: 6.25.1`: A library for routing in React applications.
  - `vite: 5.3.4`: Fast build tool.
  -`tailwindcss: 3.4.7`: Utility-first CSS framework for rapid UI development.
  - `recharts: 2.12.7`: Charting library built on React components.
  - `framer-motion: 11.3.19`: for animations and gestures in React.
  - `lucide-react: 0.417.0`: Customizable SVG icons for React.
  - [Other dependencies](https://github.com/timelygator/TimelyGator/discussions/2#discussioncomment-12105543)

- Backend:
  - `spf13/cobra` for CLI support.
  - `gorilla/mux` for API routing and middleware.
  - `gorm.io/gorm` for ORM and database management, using SQLite backend.

> Basic backend implemantaion done in `be-dev` branch.

### Discussion 4 : Documentation

For documentation, we decided to use [swaggo/swag](https://github.com/swaggo/swag) for backend API and for frontend [Vite](https://vite.dev/guide). The documentation will be hosted in the `docs/` folder and will include guides, API references, and user manuals.

> Links:
  Backend - [issue #19](https://github.com/timelygator/TimelyGator/issues/19) [draft PR #17](https://github.com/timelygator/TimelyGator/pull/17)
  Frontend - [issue #11](https://github.com/timelygator/TimelyGator/issues/11) [PR #27](https://github.com/timelygator/TimelyGator/pull/27)

### Discussion 5 : Repository settings

We discussed the repository settings and decided to:

- Enable branch protection for the `main` and `develop` branches.
- Use GitHub discussions for team communication and decision-making.
- Use GitHub issues for tracking tasks and bugs.
- `Git Flow` branching strategy will be followed with `main`, `develop`, `fe-dev`, `be-dev` branches.
- Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard for commit messages.

> Links: [discussion #3](https://github.com/timelygator/TimelyGator/discussions/3) [discussion #12](https://github.com/timelygator/TimelyGator/discussions/12)

## User stories

The user stories planned for the first sprint have been,

### Frontend

> As a potential user, I want to know how the application works and what features are available.

- [x] Create a landing page with an overview of the application.
- [ ] Add a feature list with descriptions. (reason: Basic templates have to be finalized before uploading with a feature list with description)

> As a user, I want to maintain my profile and settings.

- [x] Create a profile page with user information.
- [x] Add Theme toggle for light and dark mode. (A bit buggy ([Issue #38](https://github.com/timelygator/TimelyGator/issues/38)))
- [ ] Add settings page for customizing the application. (reason: Settings/Controls have to be finalized)

> As a user, I want a centralized dashboard for tracking time and goals.

- [x] Create a dashboard with time tracking and goal setting.
- [x] Add a daily, weekly, and monthly view for tracking time.
- [x] Show browser insights and tracking data.

### Backend

> As a developer, I want to set up the backend structure and database.

- [x] Initialize the project with Cobra CLI support.
- [x] Add GORM for ORM and database management.
- [x] Support environment files, logging, and configuration.
- [ ] Create models for user, time tracking, and goals. (reason: OAuth have to be finalized)

> As a developer, I want to set up the API routes and controllers.

- [x] Add Gorilla Mux for API routing and middleware.
- [x] Create API routes for buckets, time tracking, and goals.
- [ ] Implement CRUD operations for user data. (reason: OAuth have to be finalized)

> As a developer, I want to enable OAuth2 authentication for user login.

- [ ] Add OAuth2 support for user authentication. (reason: Google Client ID and Consent Screen requires confirmation)
- [ ] Implement login and registration endpoints. (reason: Callback URL and Redirect URL require CORS and other headers to be set)

> As a developer, I want a documentation system for API references and guides.

- [x] Set up Swaggo for API documentation.
- [x] Add godoc comments to API routes and controllers.
- [ ] Generate Swagger docs for API endpoints automatically. (reason: go generate command requires setup)

> As a developer, I want a centralized logging system for tracking errors and events.

- [x] Add logging support.
- [x] Implement logging for API requests and responses.

## Tasks

TO get a better understanding of the user stories, their status, and people contributing to them,

1. [@PulkitGarg777](https://github.com/PulkitGarg777) # Frontend
    - Setup basic `React` App using `Vite`
    - Created landing page with feature list
    - Added dashboard with time tracking and goal setting
    - Displaying daily, weekly, and monthly view for tracking time
    - Added browser insights and tracking data templates

3. [@YashDVerma](https://github.com/YashDVerma) # Frontend
    - Set up a profile page with user information
    - Implemented theme toggle for light and dark mode
    - Tested chart libraries like Apexcharts and recharts
    - Basic setup for theme

4. [@shreyansh-nayak-ufl](https://github.com/shreyansh-nayak-ufl) # Backend
    - Initialized the project and added Cobra CLI support
    - Added Gorilla Mux for API routing and middleware
    - Implemented logging support and .env configuration setup
    - Set up Swaggo for API documentation
    - Added godoc comments to API routes and controllers

5. [@siddhant-0707](https://github.com/siddhant-0707) # Backend
    - Created the Backend Architecture Design
    - Implemented GORM for ORM and database management
    - Created models for events, time tracking, and goals
    - Added API routes for user, time tracking, and goals
    - Initialize Server using Gorilla Mux

## Conclusion

The first sprint of the project focused on setting up the project structure, defining the main objectives, and evaluating feasibility. The team discussed the project structure, goals, tech stack, documentation, and repository settings. User stories were defined for frontend and backend development, and tasks were assigned to team members. The sprint laid the foundation for future development and established a roadmap for the project. The team will continue to work on the platform and provide updates in the next sprint.