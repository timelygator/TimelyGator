## SE Project : Sprint 1

This document contains the details of the first sprint of the Software Engineering project. The sprint focused on setting up the project structure, defining the main objectives, and evaluating feasibility. The goal of the sprint is to establish a solid foundation for the project and define the roadmap for future development.

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

### Discussion 1 : Project structure

For the project structure, we decided following points,
- A monorepo approach with separate folders for frontend and backend codebases.
- Documentation will be hosted in a separate folder.
- Semantic versioning will be followed for versioning the project.
- GitHub discussions will be used for team communication and decision-making, and GitHub issues will be used for tracking tasks and bugs.

The project structure is as follows:

```
project-name/
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
  - [ongoing topic](https://github.com/timelygator/TimelyGator/discussions/2)

- Backend:
  - `spf13/cobra` for CLI support.
  - `gorilla/mux` for API routing and middleware.
  - `gorm.io/gorm` for ORM and database management, using SQLite backend.

> Basic backend implemantaion done in `be-dev` branch.

### Discussion 4 : Documentation

For documentation, we decided to use [swaggo/swag](https://github.com/swaggo/swag) for backend API and for frontend [ongoing topic](https://github.com/timelygator/TimelyGator/discussions/9). The documentation will be hosted in the `docs/` folder and will include guides, API references, and user manuals.

> Links: [issue #19](https://github.com/timelygator/TimelyGator/issues/19) [draft PR #17](https://github.com/timelygator/TimelyGator/pull/17)

### Discussion 5 : Repository settings

We discussed the repository settings and decided to:
- Enable branch protection for the `main` and `develop` branches.
- Use GitHub discussions for team communication and decision-making.
- Use GitHub issues for tracking tasks and bugs.
- `Git Flow` branching strategy will be followed with `main`, `develop`, `fe-dev`, `be-dev` branches.
- Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard for commit messages.

> Links: [discussion #3](https://github.com/timelygator/TimelyGator/discussions/3) [discussion #12](https://github.com/timelygator/TimelyGator/discussions/12)
