import React, { Component } from "react";
import axios from 'axios';

const Context = React.createContext({
  isAuthenticated: false,
  authenticate: () => {},
  logout: () => {},
  user: null
});

export class LoginContextProvider extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isAuthenticated: JSON.parse(localStorage.getItem("isAuthenticated")) || false,
      authenticate: this.authenticate,
      logout: this.logout,
      user: JSON.parse(localStorage.getItem("user")) || null,
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
        this.getUser()
      })
      .catch(error => {
    	alert(error);
      })

  };

  logout = () => {
    const apiUrl = '/logout';
    axios({
      method: 'post', 
      url: `${apiUrl}`, 
    })
    this.setAuthentication(false);
    this.removeUser()

  };

  getUser = () => {
    const apiUrl = '/userInfo';
    axios({
      method: 'get', 
      url: `${apiUrl}` 
    })
      .then((response) => {
        this.setState({
          ...this.State,
          user: response.data 
        });
        localStorage.setItem("user", JSON.stringify(response.data));
      })
      .catch(error => {
    	alert(error);
      })
  };
  
  removeUser = () => {
    this.setState({
      ...this.State,
      user: null 
    })
    localStorage.removeItem("user");
  };


  setAuthentication = isAuthenticated => {
    localStorage.setItem("isAuthenticated", JSON.stringify(isAuthenticated));
    this.setState({
      ...this.State,
      isAuthenticated: isAuthenticated
    });
  };

  render() {
    const { isAuthenticated, authenticate, logout, user } = this.state;
    return (
      <Context.Provider value={{ isAuthenticated, authenticate, logout, user }}>
        {this.props.children}
      </Context.Provider>
    );
  }
}

export default {
    Provider: LoginContextProvider,
    Consumer: Context.Consumer
};

