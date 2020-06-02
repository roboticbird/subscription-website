import React from 'react';
import axios from 'axios';
import BuyPopup from '../BuyPopup';
import MessagePopup from '../MessagePopup';

class Products extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
      products: [],
      error: null,
      product_id: -1,
      showBuyPopup: false,
      showMessagePopup: false,
    }
  }

  componentDidMount() {
    const apiUrl = '/products';

    axios({method: 'get', url: `${apiUrl}`})
      .then((response) => {
        this.setState({
          products: response.data,
          isLoading: false,
        })
      })
      .catch(error => this.setState({ error:error, isLoading: false }));
  }

  toggleBuyPopup(event, id) {
    this.setState({
      showBuyPopup: !this.state.showBuyPopup,
      product_id: id
    });
  }

  toggleMessagePopup(event, message) {
    this.setState({
      showMessagePopup: !this.state.showMessagePopup,
      message: message 
    });
  }

  render() {
    const { isLoading, products } = this.state;
    return (
      <React.Fragment>
	<h1>Subscription packages</h1>
	<hr />
        {
	  !isLoading ? (
            products.map(product => {
              const { id, name, description, price, duration } = product;
              return (
                <div key={id}>
                  <h2>{name}</h2>
                  <p>{description}</p>
                  <p>Price: {(price/100).toFixed(2)} EUR, Duration: {duration} days</p>
                  {this.props.isAuthenticated ?
                   <button onClick={(e) => this.toggleBuyPopup(e, id)}>Subscribe</button>
                   : null
                  }
                  <hr />
		  {(this.state.showBuyPopup && this.state.product_id === id) ? 
		    <BuyPopup
		      product={product}
		      closeBuyPopup={(e) => this.toggleBuyPopup(e, id)}
		      toggleMessagePopup={this.toggleMessagePopup.bind(this)}
		    />
		    : null
		  }
		  {(this.state.showMessagePopup) ? 
		    <MessagePopup
                      message={this.state.message} 
		      toggleMessagePopup={this.toggleMessagePopup.bind(this)}
		    />
		    : null
		  }
                </div>
              );
            })
          ) : (
             <h3>Loading...</h3>
          )
        }
      </React.Fragment> 
    );
  }

}

export default Products;
