# React UI Development Guide

This directory contains the modernized React-based user interface for Archive Proxy.

## Tech Stack

- **React 19** - UI framework
- **Vite** - Build tool and dev server
- **Ant Design 6** - UI component library
- **StreamSaver.js** - Client-side file streaming

## Development

### Prerequisites

- Node.js 18+ and npm

### Setup

```bash
cd web/react-ui
npm install
```

### Development Server

Start the development server with hot reload:

```bash
npm run dev
```

The dev server runs on `http://localhost:5173` and proxies API requests to `http://localhost:8080`.

**Note:** You need to run the Archive Proxy Go server separately on port 8080 for the API endpoints to work.

### Building for Production

Build the optimized production bundle:

```bash
npm run build
```

The built files are output to `web/dist/` which is embedded into the Go binary.

### Linting

```bash
npm run lint
```

## Project Structure

```
src/
  ├── App.jsx          # Main application component
  ├── App.css          # Application styles
  ├── main.jsx         # React entry point
  └── index.css        # Global styles
public/
  └── StreamSaver.js   # File streaming library
```

## Features

- **URL Input** - Enter archive URL with validation
- **Encoding Selection** - Support for 40+ character encodings
- **File Listing** - Fetch and display archive contents
- **Multi-Select** - Select multiple files for download
- **Download** - Stream selected files as a zip
- **Error Handling** - User-friendly error messages
- **Responsive Design** - Works on all screen sizes

## API Integration

The UI integrates with the following Archive Proxy endpoints:

- `GET /list?url=<url>&charset=<charset>` - List archive contents
- `POST /pack?url=<url>&charset=<charset>` - Download selected files

## Contributing

When making changes to the UI:

1. Make your changes in the `src/` directory
2. Test in development mode with `npm run dev`
3. Build for production with `npm run build`
4. Verify the built files in `web/dist/`
5. Rebuild the Go server to embed the new UI
