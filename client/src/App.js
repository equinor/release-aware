import React from 'react';
import styled from 'styled-components';
import ReleasesContainer from './ReleasesContainer';
import Spinner from './Spinner';
import axios from 'axios';

const Header = styled.h3`
  margin: 0px;
  padding: 10px;
  color: white;
`;

const AppContainer = styled.div`  
  margin: 0px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 20px;
`;

const P = styled.p`
  margin: 0px;
`;

function CouldNotFetch({lastSuccessfulFetch}) {
    let lastFetch;
    if (lastSuccessfulFetch) {
        const temp_string = lastSuccessfulFetch.toString();
        lastFetch = temp_string
            .split(' ')
            .splice(0, 5)
            .join(' ');
    } else {
        lastFetch = 'none';
    }
    return (
        <>
            <P>Could not fetch data</P>
            <P>Last updated; {lastFetch}</P>
        </>
    );
}

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            releases: []
        };
    }

    refetchData() {
        axios
            .all([
                axios.get("/api/releases"),
                axios.get("/api/helmreleases")
            ])
          .then(axios.spread((...responses) => {
            this.setState({ 
                isLoaded: true,
                releases: responses[0].data.concat(responses[1].data),
                error: null,
             })
          }))
            // Note: it's important to handle errors here
            // instead of a catch() block so that we don't swallow
            // exceptions from actual bugs in components.
            .catch(error => {
                this.setState({
                    isLoaded: true,
                    error,
                    lastSuccessfulFetch: new Date(),
                });
            });
    }

    componentDidMount() {
        this.refetchData();
        this.interval = setInterval(() => this.refetchData(), 100000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    render() {
        const {
            error,
            isLoaded,
            releases,
            lastSuccessfulFetch
        } = this.state;

        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <Spinner isLoading={true}/>;
        } else {
            return (
                <AppContainer>
                    <Header>Releases</Header>
                    {error && <CouldNotFetch lastSuccessfulFetch={lastSuccessfulFetch}/>}
                    <ReleasesContainer releases={releases}/>
                </AppContainer>
            );
        }
    }
}

export default App;
