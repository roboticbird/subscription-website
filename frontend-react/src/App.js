import React from 'react';
import Login from './components/Login';
import LoginContext from "./components/LoginContext";
import Navigation from './components/Navigation';
import Products from './components/Products';
import Subscriptions from './components/Subscriptions';
import { Router, Redirect } from "@reach/router"
import './App.css';

export default function App() {
  let NavBar = () => (
      <LoginContext.Consumer>
        {props => <Navigation 
          isAuthenticated={props.isAuthenticated} 
          logout={props.logout}
          user={props.user} 
        />}
      </LoginContext.Consumer>
  );

  let ProductPage = () => (
    <div>
      <NavBar />
      <LoginContext.Consumer>
        {props => <Products isAuthenticated={props.isAuthenticated} />}
      </LoginContext.Consumer>
    </div>
  );
  
  let SubscriptionPage = () => (
    <div>
      <NavBar />
      <LoginContext.Consumer>
        {props => <Subscriptions isAuthenticated={props.isAuthenticated} />}
      </LoginContext.Consumer>
    </div>
  );

  let LoginPage = () => (  
    <div>
      <NavBar />
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
       <ProductPage path="/" />
       <SubscriptionPage path="subscriptions" />
       <LoginPage path="login" />
      </Router>
    </LoginContext.Provider>
    </div>
  );
}

