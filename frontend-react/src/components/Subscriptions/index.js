import React from 'react';
import axios from 'axios';
import StatusPopup from '../StatusPopup';

class Subscriptions extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
      subscriptions: [],
      error: null,
      subscriptionId: -1,
      reqAction: "",
    }

    this.loadSubscriptions = this.loadSubscriptions.bind(this)
  }

  componentDidMount() {
    console.log("componentDidMount")
    this.loadSubscriptions();
    this.setState({
          showPopup: false,
    })
  }

  togglePopup(event, id, ACTION) {
    console.log("togglePopup")
    this.setState({
      showPopup: !this.state.showPopup,
      subscriptionId: id,
      reqAction: ACTION
    });
    if (!this.state.showPopup) {
      console.log("reload")
      this.loadSubscriptions();
    }
  }
  
  loadSubscriptions() {
    this.setState({
      isLoading: true,
    })
    const apiUrl = '/subscriptions';
    axios({method: 'get', url: `${apiUrl}`})
      .then((response) => {
        this.setState({
          subscriptions: response.data,
          isLoading: false,
        })
      })
      .catch(error => this.setState({ error, isLoading: false }));
  }

  render() {
    const { isLoading, subscriptions } = this.state;
    return (
      <React.Fragment>
	<h1>My subscription packages</h1>
	<hr />
        {
	  !isLoading ? (
            subscriptions == null ? <></> : (
            subscriptions.map(subscription => {

              const { id, product, startDate, subStatus, nextStatus } = subscription 
              const { name, description, price, duration } = product;
              return (
                <div key={id}>
                  <h2>{name}</h2>
                  <p>{description}</p>
                  <p>Price: {(price/100).toFixed(2)} EUR, Duration: {duration} days</p>
                  <p>Start date: {new Date(startDate).toLocaleDateString("en-GB")}, Status: {subStatus}</p>
                  <p>{nextStatus === "PAUSED" ? "You have paused your next payment" : "Your next month is currently queued"}</p>
                  <button onClick={(e) => this.togglePopup(e, id, nextStatus === "QUEUED" ? "PAUSE" : "QUEUE")}>
                  {nextStatus === "QUEUED" ? "Pause" : "Start" }
                  </button>
                  <button onClick={(e) => this.togglePopup(e, id, "CANCEL")}>
                  Cancel 
                  </button>
                  <hr />
		  {(this.state.showPopup && this.state.subscriptionId === id) ? 
		    <StatusPopup
		      product={product}
                      action={this.state.reqAction}
		      closePopup={(e) => this.togglePopup(e, product.id)}
		    />
		    : null
		  }
                </div>
              );
            })
          )) : (
             <h3>Loading...</h3>
        )}
      </React.Fragment> 
    );
  }

}

export default Subscriptions;
