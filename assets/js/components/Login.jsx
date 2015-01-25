var React = require('react');
var Router = require('react-router');
var Actions = require('../Actions');
var UserStore = require('../stores/UserStore');


var Login = React.createClass({

    mixins: [Router.Navigation],

    handleSubmit: function(event) {
        event.preventDefault();
        var identifier = this.refs.identifier.getDOMNode().value;
        var password = this.refs.password.getDOMNode().value;
        Actions.login(identifier, password);
    },

    componentWillMount: function() {
        UserStore.addChangeListener(this._onChange);
    },

    componentWillUnmount: function () {
        UserStore.removeChangeListener(this._onChange);
    },

    render: function() {
        return (
        <form role="form" name="form" onSubmit={this.handleSubmit}>
            <div className="form-group">
                <input type="text" ref="identifier" required="required" className="form-control" placeholder="Name or email address" />
            </div>

            <div className="form-group">
                <input type="password" ref="password" required="required" className="form-control" placeholder="Password" />
            </div>

            <button type="submit" className="btn btn-primary">Login</button><br />
            <a href="#/recoverpass">Forgot your password?</a>
        </form>
        );
    },

    _onChange: function () {
        user = UserStore.getUser();
        if (user) {
            this.transitionTo("popular");
        }
    }

});

module.exports = Login;
