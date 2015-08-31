import React, { PropTypes } from 'react';
import _ from 'lodash';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input, 
         ButtonInput 
        } from 'react-bootstrap';

import * as ActionCreators from '../actions';


export default class Upload extends React.Component {
  render() {
    return (
      <div className="row">
          <div className="col-md-6">
              <form name="form" role="form" enctype="multipart/form-data">
                <Input name="title" type="text" ref="title" label="Title" />
                <Input name="tags" type="text" ref="tags" label="Tags" placeholder="Separate with spaces" />
                <Input name="tags" type="file" ref="photo" label="Photo" />
                <ButtonInput bsStyle="primary">Upload</ButtonInput>
              </form>
          </div>
          <div className="col-md-6">
              <div className="thumbnail">
                  <img src="" />
              </div>
          </div>
      </div>
    );
  }
}
