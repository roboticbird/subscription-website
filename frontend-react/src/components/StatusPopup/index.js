import React from "react"
import axios from 'axios';
import './styles.css';

class StatusPopup extends React.Component {
  constructor (props) {
    super(props)

    this.handleSubmit = this.handleSubmit.bind(this)
  }

  handleSubmit(event) {
    const apiUrl = '/updateSubscription';
    
    axios({
      method: 'post', 
      url: `${apiUrl}`, 
      data: {
        productId: this.props.product.id,
        Action: this.props.action
      }
    })
      .then((response) => {
        this.props.closePopup();
        console.log("sucess")
      })
      .catch(error => {
        this.props.closePopup();
    	console.log(error);
      })


  }

  render() {
    return (
      <div className='StatusPopup'>
        <div className='StatusPopupInner'>
          {
            {
              "QUEUE" : (<h1>Activate: {this.props.product.name}</h1>),
              "PAUSE" : (<h1>Pause: {this.props.product.name}</h1>),
              "CANCEL" : (<h1>Cancel: {this.props.product.name}</h1>),
            }[this.props.action]
          }
          <p>Price: {(this.props.product.price/100).toFixed(2)} EUR, Duration: {this.props.product.duration} days</p>
          <button onClick={(e) => this.handleSubmit(e)}>Confirm</button>
        </div>
      </div>
    );
  }
}

export default StatusPopup;
