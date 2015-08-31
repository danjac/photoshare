import React from 'react';


export class Facon extends React.Component {

  static propTypes = {
    name: React.PropTypes.string.isRequired
  }

  render() {
    const className = `fa fa-${this.props.name}`;
    return (
      <i className={className}></i>
    );
  }

}
