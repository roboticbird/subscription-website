import React, { Component } from 'react';
import Products from './components/Products';
import Login from './components/Login';
import './App.css';


export const GlobalStateContext = React.createContext();


class App extends Component {
  setAuthentication = isAuthenticated => {
    this.setAuthentication({ isAuthenticated });
  };
  
  globalState = {
    isAuthenticated: false,
    setAuthentication: () => {},
  };


  render() {
    return (
      <div className="App">
        <GlobalStateContext.Provider value={this.globalState}>
          <Login />
          <Products />
        </GlobalStateContext.Provider>
      </div>
    );
  }
}

export default App;
