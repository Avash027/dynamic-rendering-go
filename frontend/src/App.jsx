import { useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="App">
      <h1>Counter</h1>
      <p>Count: {count}</p>
      <button onClick={() => setCount(count + 1)}>Increment</button>
      <button onClick={() => setCount(count - 1)}>Decrement</button>

      <h1>Todo List</h1>
      <ul>
        <li>Learn React</li>
        <li>Learn Redux</li>
        <li>Learn React-Redux</li>
      </ul>

    </div>
  )
}

export default App
