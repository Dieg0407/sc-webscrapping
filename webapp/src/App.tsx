import "./App.css";
import Header from "./components/Header/Header";
import OperationController from "./components/OperationController/OperationController";

function App() {
  return (
    <>
      <Header text="Sitio de administración de reportes del SEACE" />
      <OperationController />
    </>
  );
}

export default App;
