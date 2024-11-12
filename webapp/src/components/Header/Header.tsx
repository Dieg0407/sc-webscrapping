import styles from "./Header.module.css";

export type HeaderProps = {
  /**
   * The text to display in the header
   */
  text: string;
};

function Header({ text }: HeaderProps) {
  return (
    <header className={styles.header}>
      <h1>{text}</h1>
    </header>
  );
}

export default Header;
