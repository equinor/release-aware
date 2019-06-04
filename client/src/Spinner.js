import React from 'react';
import styled from 'styled-components';

const Wrapper = styled.div` 
  .loader {
      margin: 20px auto;
      border: 8px solid #bbbbbb;
      border-top: 8px solid #000;
      border-radius: 50%;
      width: 160px;
      height: 160px;
      animation: spin 1.5s linear infinite;
  }
  
  @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
  }
`;

const Spinner = (props) => {
    const {isLoading, component = null} = props;
    if (isLoading) {
        return (
            <Wrapper>
                <div className="loader" />
            </Wrapper>
        );
    }

    return (
        <div>
            {component}
            {props.children}
        </div>
    );
};


export default Spinner;