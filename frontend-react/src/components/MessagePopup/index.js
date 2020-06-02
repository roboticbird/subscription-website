import React from "react"
import './styles.css';

class MessagePopup extends React.Component {
  render() {
    return (
      <div className='MessagePopup'>
        <div className='MessagePopupInner'>
          <p>{this.props.message}</p>
          <button onClick={(e) => this.props.toggleMessagePopup(e,"")}>Ok</button>
        </div>
      </div>
    );
  }
}

export default MessagePopup;
