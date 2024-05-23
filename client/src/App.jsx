import './App.css'
import {useState, useEffect} from 'react'

function App() {
  const [word, setWord] = useState('')
  const [translation, setTranslation] = useState('')

  async function loadWords(){
    const response = await fetch(import.meta.env.VITE_API + '/palabras')
    const data = await response.json()
    console.log(data)
  }

  useEffect(() => {
    loadWords()
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    const response = await fetch(import.meta.env.VITE_API + '/palabra', {
      method: 'POST',
      body: JSON.stringify({WordText: word, Translation: translation}),
      headers: {
        'Content-Type': 'application/json'
      }
    })
    const data = await response.json()
    console.log(data)
    loadWords()
  }

  return (
    <>
      <div className="title-header">
        <p>Traduccion de español a Q'eqchi'</p>
      </div>
      <div className="body-container">
        <form onSubmit={handleSubmit}>
          <input type="word" placeholder="Ingresa una palabra" onChange={e => setWord(e.target.value)}/>
          <br />
          <input type="translation" placeholder="Ingresa la traducción" onChange={e => setTranslation(e.target.value)}/>
          <br />
          <button>Guardar</button>
        </form>
      </div>
    </>
  )
}

export default App
