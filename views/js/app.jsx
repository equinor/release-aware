class HelloWorld extends React.Component {
    render() {
      return (
        <div className="container">
          <div className="row">
            <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
              <h1>go-gin-react-helloworld</h1>
              <h2>Hello World!</h2>
              <h3>Hello World web app created with Golang + Gin + React</h3>
              <h4><a href="https://github.com/cskonopka/go-gin-react">Github</a></h4>
            </div>
          </div>
        </div>
      );
    }
  }
  
  ReactDOM.render(<HelloWorld />, document.getElementById("app"));