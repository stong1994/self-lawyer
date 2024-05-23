import React, { useState } from "react";

function App() {
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");

  const handleSubmit = (event) => {
    event.preventDefault();
    // Here you would add your logic for generating an answer
    // visit http://localhost:8888/chatbot to see the answer
    // Define the question
    let questionData = { question: question };

    // Make a POST request
    fetch("http://localhost:8888/chat", {
      method: "POST", // Specify the method
      headers: {
        "Content-Type": "application/json", // Set the content type to JSON
      },
      body: JSON.stringify(questionData), // Convert the question to a JSON string
    })
      .then((response) => response.json()) // Parse the response as JSON
      .then((data) => {
        console.log(data);
        setAnswer(data.answer);
      }) // Log the data to the console
      .catch((error) => console.error("Error:", error)); // Log any errors
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <label>
          Ask a question:
          <input
            type="text"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
          />
        </label>
        <input type="submit" value="Submit" />
      </form>
      <p>{answer}</p>
    </div>
  );
}

export default App;
