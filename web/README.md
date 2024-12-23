## Installation and Usage

### Clone the Project

First, clone this project to your local machine:

```bash
git clone https://github.com/GabbyWorld/all-time-high
cd all-time-high
```

### Install Dependencies
Install the required dependencies:

```bash
npm install
```

Or use Yarn:
```bash
yarn install
```

### Run the Project
Start the development server locally:
```bash
npm run dev
```

Or use Yarn:
```bash
yarn dev
```

Then open your browser and visit http://localhost:3000 to view the project.

## Project Structure
Here is the basic directory structure of the project:

```bash
├── public/                  # Public files, static resources
│   └── index.html           # Project's HTML template
├── src/                     # Source code folder
│   ├── assets/              # Images, fonts, and other resources
│   ├── components/          # React components
│   ├── hooks/               # Custom hooks
│   ├── pages/               # Page components
│   ├── App.tsx              # The root component of the application
│   ├── index.tsx            # The entry file of the project
├── .gitignore               # Git ignore file
├── package.json             # Project's configuration and dependencies
├── README.md                # Project's documentation
└── tsconfig.json            # TypeScript configuration file

```
## Tech Stack
- React: For building the user interface
- TypeScript: For static type checking
- Redux: For state management
- Chakra UI: UI component library
- Vite: Frontend build tool
- Axios: For handling HTTP requests
- react-router-dom: For routing management
