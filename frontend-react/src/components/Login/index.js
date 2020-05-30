import React from 'react';
import { Button, FormGroup, FormControl, FormLabel } from "react-bootstrap";
import axios from 'axios';



class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
	email: "",
	password: "",
	
    };
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  validateForm() {
    const { email, password } = this.state;
    return email.length > 0 && password.length > 0;
  }

  handleSubmit(event) {
    const apiUrl = '/login';
    const { email, password } = this.state;
    event.preventDefault();

    axios({
      method: 'post', 
      url: `${apiUrl}`, 
      data: {
        email: email,
        password: password
      }
    })
      .then((response) => {
        localStorage.setItem("isAuthenticated", true);
        this.props.setAuthentication(true);
      })
      .catch(error => {
    	alert(error);
      })
  }

  render() {
      const { email, password } = this.state;

      return (<div className="Login">
	<form onSubmit={this.handleSubmit}>
	  <FormGroup controlId="email" bssize="large">
	    <FormLabel>Email</FormLabel>
	    <FormControl
	      autoFocus
	      type="email"
	      value={email}
	      onChange={e => this.setState({email: e.target.value})}
	    />
	  </FormGroup>
	  <FormGroup controlId="password" bssize="large">
	    <FormLabel>Password</FormLabel>
	    <FormControl
	      value={password}
	      onChange={e => this.setState({password: e.target.value})}
	      type="password"
	    />
	  </FormGroup>
	  <Button block bssize="large" disabled={!this.validateForm()} type="submit">
	    Login
	  </Button>
	</form>
      </div>);
  }
}

export default Login
