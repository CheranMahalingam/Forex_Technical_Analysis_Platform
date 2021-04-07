import React from "react";
import NavigationMenu from "../components/NavigationMenu/NavigationMenu";
import logo from "../logo.svg";
import "../App.css";

function HomePage() {
  return (
    <div
      style={{
        backgroundColor: "#000000",
        minHeight: "100%",
        height: "auto",
        top: 0,
        left: 0,
        width: "100%",
        position: "absolute",
      }}
    >
      <NavigationMenu />
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/App.js</code> and save to reload.
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
        </header>
      </div>
    </div>
  );
}

export default HomePage;
