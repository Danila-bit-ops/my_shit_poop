function App() {
  const [paramID, setParamID] = useState('');
  const [responseData, setResponseData] = useState(null);

  const handleButtonClick = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/get-by-param-id?paramID=${paramID}`);
      const data = await response.json();
      setResponseData(data);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  return (
    <div>
      <label>
        Param ID:
        <input type="text" value={paramID} onChange={(e) => setParamID(e.target.value)} />
      </label>
      <button onClick={handleButtonClick}>Send GET Request</button>

      {responseData && (
        <div>
          <h2>Response:</h2>
          <pre>{JSON.stringify(responseData, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}

export default App;