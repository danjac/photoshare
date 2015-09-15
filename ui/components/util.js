import React from 'react';


export class Loader extends React.Component {
  render() {
    return (
      <div className="col-md-6 col-md-offset-6">
        <img src="/img/ajax-loader.gif" alt="" />
      </div>
    );
  }
}

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
