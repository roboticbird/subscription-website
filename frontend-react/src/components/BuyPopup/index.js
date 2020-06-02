import React from "react"
import axios from 'axios';
import './styles.css';

class BuyPopup extends React.Component {
  constructor (props) {
    super(props)

    this.handleSubmit = this.handleSubmit.bind(this)
  }

  handleSubmit(event) {
    const apiUrl = '/newSubscription';
    axios({
      method: 'post', 
      url: `${apiUrl}`, 
      data: {
        productId: this.props.product.id
      }
    })
      .then((response) => {
        this.props.toggleMessagePopup(null, "Subscription successful."); 
        this.props.closeBuyPopup();
      })
      .catch(error => {
        var errMsg = "Error purchasing subscription"
        if (error.response.status === 304) {
          errMsg = "You are already subscribed"
        }
        this.props.toggleMessagePopup(null, errMsg); 
        this.props.closeBuyPopup();
      })


  }

  render() {
    return (
      <div className='BuyPopup'>
        <div className='BuyPopupInner'>
          <h1>Buy: {this.props.product.name}</h1>
          <p>Price: {(this.props.product.price/100).toFixed(2)} EUR, Duration: {this.props.product.duration} days</p>
          <button onClick={(e) => this.handleSubmit(e)}>Confirm</button>
        </div>
      </div>
    );
  }
}

export default BuyPopup;
