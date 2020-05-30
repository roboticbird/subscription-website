import React from 'react';

class Navigation extends React.Component {
  render() {
    console.log(this.props);
    return(
      <div className="Navigation">
        <ul id="nav">
          <li><a href="/">Home</a></li>
          {this.props.isAuthenticated ? 
            this.props.user == null ? (
              <></>
            ): (
              <li>User: {this.props.user.name}</li>
            )
           : (
            <li><a href="login">Login</a></li>
          )}
        </ul>
      </div>
    );
  }
}
export default Navigation;
