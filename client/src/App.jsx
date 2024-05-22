import './App.css'

function App() {
  return (
    <>
      <div className="title-header">
        <p>Traduccion de español a Q'eqchi'</p>
      </div>
      <div className="body-container">
        <button onClick={async () => {
          const response = await fetch('/users')
          const data = await response.json()
          console.log(data)
        }}>Obtener Datos</button>
      </div>
    </>
  )
}

export default App
