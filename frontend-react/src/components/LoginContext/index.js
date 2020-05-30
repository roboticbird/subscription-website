import React, { Component } from "react";
import axios from 'axios';

const Context = React.createContext({
  isAuthenticated: false,
  authenticate: () => {}
});

export class LoginContextProvider extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isAuthenticated: JSON.parse(localStorage.getItem("isAuthenticated")) || false,
      authenticate: this.authenticate
    };
  }

  authenticate = (email, password) => {
    const apiUrl = '/login';
    axios({
      method: 'post', 
      url: `${apiUrl}`, 
      data: {
        email: email,
        password: password
      }
    })
      .then((response) => {
        this.setAuthentication(true);
      })
      .catch(error => {
    	alert(error);
      })

  };

  setAuthentication = isAuthenticated => {
    localStorage.setItem("isAuthenticated", JSON.stringify(isAuthenticated));
    this.setState({
      ...this.State,
      isAuthenticated: isAuthenticated
    });
  };

  render() {
    const { isAuthenticated, authenticate } = this.state;
    return (
      <Context.Provider value={{ isAuthenticated, authenticate }}>
        {this.props.children}
      </Context.Provider>
    );
  }
}

export default {
    Provider: LoginContextProvider,
    Consumer: Context.Consumer
};

