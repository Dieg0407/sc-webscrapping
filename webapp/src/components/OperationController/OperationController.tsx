import { useState } from "react";
import styles from "./OperationController.module.css";

function OperationController() {
  const [selected, setSelected] = useState("explore");

  return (
    <div className={styles.container}>
      <button
        className={`${styles.button} ${
          selected === "explore" ? styles.active : styles.deactivated
        }`}
        onClick={() => setSelected("explore")}
      >
        Explorar Reportes
      </button>
      <button
        className={`${styles.button} ${
          selected === "generate" ? styles.active : styles.deactivated
        }`}
        onClick={() => setSelected("generate")}
      >
        Generar Reporte
      </button>
    </div>
  );
}

export default OperationController;
