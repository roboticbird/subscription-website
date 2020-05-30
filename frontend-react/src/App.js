import React from 'react';
import Products from './components/Products';
import Navigation from './components/Navigation';
import Login from './components/Login';
import LoginContext from "./components/LoginContext";
import { Router, Redirect } from "@reach/router"
import './App.css';

export default function App() {
  let Home = () => (
    <div>
      <LoginContext.Consumer>
        {props => <Navigation isAuthenticated={props.isAuthenticated} user={props.user} />}
      </LoginContext.Consumer>
      <Products />
    </div>
  );

  let LoginPage = () => (  
    <div>
      <LoginContext.Consumer>
        {props => <Navigation isAuthenticated={props.isAuthenticated} user={props.user} />}
      </LoginContext.Consumer>
      <LoginContext.Consumer>
        {props =>
          !props.isAuthenticated ? (
            <Login authenticate={props.authenticate} />
          ) : (
            <Redirect to="/" noThrow />
          )
        }
      </LoginContext.Consumer>
    </div>
  );

  return (
    <div className="App">
    <LoginContext.Provider>
      <Router>
       <Home path="/" />
       <LoginPage path="login" />
      </Router>
    </LoginContext.Provider>
    </div>
  );
}

