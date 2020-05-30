import React from 'react';
import Products from './components/Products';
import Login from './components/Login';
import LoginContext from "./components/LoginContext";
import './App.css';

export default function App() {
  return (
    <div className="App">
      <LoginContext.Provider>
        <LoginContext.Consumer>
          {props =>
            !props.isAuthenticated ? (
              <Login setAuthentication={props.setAuthentication} />
            ) : (
              <div>Logged in</div>
            )
          }
        </LoginContext.Consumer>
      </LoginContext.Provider>
      <Products />
    </div>
  );
}

