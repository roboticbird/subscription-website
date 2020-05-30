import React, { Component } from "react";

const Context = React.createContext({
  isAuthenticated: false,
  setAuthentication: () => {}
});

export class LoginContextProvider extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isAuthenticated: false,
      setAuthentication: this.setAuthentication
    };
  }

  setAuthentication = isAuthenticated => {
    this.setState({
      ...this.State,
      isAuthenticated: isAuthenticated
    });
  };

  render() {
    const { isAuthenticated, setAuthentication } = this.state;
    return (
      <Context.Provider value={{ isAuthenticated, setAuthentication }}>
        {this.props.children}
      </Context.Provider>
    );
  }
}

export default {
    Provider: LoginContextProvider,
    Consumer: Context.Consumer
};

