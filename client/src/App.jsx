import './App.css'
import {useState, useEffect} from 'react'
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function App() {
  //VARIABLES PARA INSERTAR PALABRAS
  const [word, setWord] = useState('')
  const [translation, setTranslation] = useState('')
  //PALABRA QUE SE BUSCARA
  const [wordForSearch, setWordForSearch] = useState('')
  //PALABRA TRADUCIDA
  const [wordTranslated, setWordTranslated] = useState('')

  const [similarWords, setSimilarWords] = useState([]);

  const [voicesSpeech, setVoices] = useState([]);

  useEffect(() => {
    // Función para cargar las voces disponibles
    const loadVoices = () => {
      const synth = window.speechSynthesis;
      const voices = synth.getVoices();
      setVoices(voices);
    };

    loadVoices();
  }, []); 

  //SE LE ASIGNA EL TEXTO ESCRITO A LA VARIABLE
  const handleTextAreaChange = (e) => {
    setWordForSearch(e.target.value);
    handleGetSuggestions(e.target.value);
  }

  const setSuggestionText = (text) => {
    setWordForSearch(text)
  }

  //PETICION PARA BUSCAR UNA PALABRA
  const handleTranslate = async (e) => {
    e.preventDefault()

    if (wordForSearch === '') {
      toast.warning('Ingrese una palabra')
      return
    }

    const response = await fetch(import.meta.env.VITE_API + `/buscar?palabra=${encodeURIComponent(wordForSearch)}`);
    const data = await response.json();
    if(data.statusCode == 404) {
      toast.error(data.error);
      return;
    }
    setWordTranslated(data.Translation)
  }

  const handleGetSuggestions = async (word) => {
    try {
      const response = await fetch(import.meta.env.VITE_API + `/sugerencias?palabra=${encodeURIComponent(word)}`);
      if (!response.ok) {
        throw new Error('Error en la solicitud: ' + response.statusText);
      }
      const data = await response.json();
      
      setSimilarWords(data);
    } catch (error) {
      toast.error('Hubo un problema con la solicitud:', error);
      toast.error('Hubo un problema al obtener palabras similares');
    }
  }

  //PETICION PARA INSERTAR UNA PALABRA A BD
  const handleSubmit = async (e) => {
    e.preventDefault()
    const response = await fetch(import.meta.env.VITE_API + '/palabra', {
      method: 'POST',
      body: JSON.stringify({WordText: word, Translation: translation}),
      headers: {
        'Content-Type': 'application/json'
      }
    })
    await response.json()
  }

  const handleEraseAll = (e) => {
    e.preventDefault()

    setWordForSearch('');
    setWordTranslated('');
    toast.success('Campos limpiados')
  }

  const handleRead = (e) => {
    e.preventDefault();

    const synth = window.speechSynthesis;
    if (synth.speaking) {
      toast.error('Ya está hablando.');
      return;
    }
    if (wordTranslated !== '') {
      const voices = synth.getVoices()
      const utterance = new SpeechSynthesisUtterance(wordTranslated);
      utterance.voice = voices[15]; //2, 6, 7, 8, 9*, 13, 14, 15*, 18,
      utterance.onend = () => {
        toast.done('Narración terminada.');
      };
      utterance.onerror = (event) => {
        toast.error('Ocurrió un error durante la narración:', event);
      };
      synth.speak(utterance);
    }
  }

  return (
    <>
      <div className="title-header">
        <p>Traduccion de español a Q'eqchi'</p>
      </div>
      <div className="body-container">
        <form>
          <div className='container'>
            <div className="textarea-container">
              <h1>Español</h1>
              <textarea value={wordForSearch} onChange={handleTextAreaChange}></textarea>
              <div className='buttons-container'>
                <button onClick={handleTranslate}>Traducir</button>
                <button onClick={handleEraseAll}>Borrar</button>
              </div>
            </div>
            <div className="textarea-container">
              <h1>Q'eqchi'</h1>
              <textarea readOnly value={wordTranslated}></textarea>
              <div className='buttons-container'>
                <button onClick={handleRead}>Voz</button>
              </div>
            </div>
          </div>
          <div className='suggestion-container'>
            <div className="suggested-textarea">
              <h1>Palabras Sugeridas</h1>
              <div className='txt'>
                <ol>
                  {similarWords.map((word, index) => (
                    <li key={index}>
                      <p className='suggestion' onClick={() => setSuggestionText(word.WordText)}>{word.WordText}</p>
                    </li>
                  ))}
                </ol>
              </div>
            </div>
          </div>
          {/* <input type="word" placeholder="Ingresa una palabra" onChange={e => setWord(e.target.value)}/>
          <br />
          <input type="translation" placeholder="Ingresa la traducción" onChange={e => setTranslation(e.target.value)}/>
          <br />
          <button onClick={handleSubmit}>Guardar</button> */}
        </form>
      </div>
      <ToastContainer />
    </>
  )
}

export default App
