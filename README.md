# All Time High - Go Web Backend & React Frontend

Welcome to **All Time High**, a cutting-edge web application designed to provide a seamless experience for users interacting with the Solana blockchain. This repository showcases a robust Go-based backend and a dynamic React frontend, making it perfect for your next hackathon!

## 🚀 Project Overview

### Backend
- **Built with Go**: Utilizing the Gin framework for high performance.
- **Data Management**: PostgreSQL for data storage, with GORM as the ORM for efficient database interactions.
- **Configuration**: Managed by Viper for easy environment variable handling.
- **User Authentication**: Integrates Phantom wallet for secure user authentication and manages sessions with JWT.
- **API Documentation**: Auto-generated Swagger docs for easy endpoint reference and testing.

### Frontend
- **React & TypeScript**: A modern frontend built with React, bundled by Vite for fast development.
- **UI Components**: Utilizes Chakra UI for a responsive and accessible user interface.
- **Seamless Integration**: Facilitates Phantom wallet integration for a smooth Solana-based experience.

## 📋 Prerequisites

To get started, ensure you have the following installed:
- Go 1.23 or later
- Node.js and npm (or yarn/pnpm) for the React frontend
- A running PostgreSQL instance
- (Optional) Docker/Docker Compose for container-based deployment

## 🏁 Getting Started

Follow these steps to set up the project locally:

1. **Clone the repository**:
   ```bash
   git clone https://github.com/GabbyWorld/all-time-high-backend.git
   cd all-time-high-backend
   ```

2. **Configure environment variables**:
   - In the server directory, copy the sample file:
     ```bash
     cp .env.example .env
     ```
     Fill in the required fields (database credentials, OpenAI keys, etc.).
   - For the frontend, set environment variables within Vite configs or .env files as necessary.

3. **Install dependencies**:
   - **Backend**:
     ```bash
     make deps
     ```
   - **Frontend** (inside the "web" folder):
     ```bash
     npm install
     ```
     (You may use yarn or pnpm if you prefer.)

4. **Build and run the Go server**:
   ```bash
   make run
   ```
   By default, the server listens on port 9100 (configurable in your .env file).

5. **Generate Swagger documentation** (optional):
   If you want to regenerate the docs:
   ```bash
   make swagger
   ```

6. **Run the frontend**:
   - From the "web" directory:
     ```bash
     npm run dev
     ```
   - The default development port is 3000. Adjust in vite.config.ts if needed.

## 🐳 Docker Usage

A sample Dockerfile is provided for the Go backend. You can build and run it as follows:

```bash
# Add Docker instructions here
```

## 🎉 Contributing

We welcome contributions! Feel free to submit issues or pull requests to help improve the project.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Thank you for checking out **All Time High**! We hope you enjoy using it as much as we enjoyed building it. Good luck at the hackathon!
