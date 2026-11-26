import { UploadPage } from './pages/UploadPage'

export default function App() {
  return (
    <div className="app">
      <header className="app__header">
        <span className="brand">
          Resume<span className="brand__accent">Builder</span>
        </span>
      </header>

      <main className="app__main">
        <UploadPage />
      </main>

      <footer className="app__footer">Built with React and Go.</footer>
    </div>
  )
}
