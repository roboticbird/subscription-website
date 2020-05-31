import React from 'react';

class Navigation extends React.Component {
  render() {
    return(
      <div className="Navigation">
        <ul id="nav">
          <li><a href="/">Products</a></li>
          {this.props.isAuthenticated ? 
            this.props.user == null ? (
              <></>
            ): (
              <>
              <li><a href="subscriptions">Subscriptions</a></li>
              <li>User: {this.props.user.name}</li>
              <li><button onClick={() => this.props.logout()}>Logout</button></li>
              </>
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
