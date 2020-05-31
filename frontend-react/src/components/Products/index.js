import React from 'react';
import axios from 'axios';
import BuyPopup from '../BuyPopup';

class Products extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
      products: [],
      error: null,
      product_id: -1,
    }
  }

  componentDidMount() {
    const apiUrl = '/products';

    axios({method: 'get', url: `${apiUrl}`})
      .then((response) => {
        this.setState({
          products: response.data,
          isLoading: false,
          showPopup: false,
        })
      })
      .catch(error => this.setState({ error, isLoading: false }));
  }

  togglePopup(event, id) {
    this.setState({
      showPopup: !this.state.showPopup,
      product_id: id
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
                   <button onClick={(e) => this.togglePopup(e, id)}>Subscribe</button>
                   : null
                  }
                  <hr />
		  {(this.state.showPopup && this.state.product_id === id) ? 
		    <BuyPopup
		      product={product}
		      closePopup={(e) => this.togglePopup(e, id)}
		    />
		    : null
		  }
                </div>
              );
            })
          ) : (
             <h3>Loading...</h3>
        )}
      </React.Fragment> 
    );
  }

}

export default Products;
