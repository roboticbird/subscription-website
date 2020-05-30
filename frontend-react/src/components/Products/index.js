import React from 'react';
import axios from 'axios';

class Products extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isLoading: true,
      products: [],
      error: null
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
      .catch(error => this.setState({ error, isLoading: false }));
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
                  <p>Price: {price}, Duration: {duration}</p>
                  <hr />
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
