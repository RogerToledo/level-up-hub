// frontend/src/app/page.tsx
import { redirect } from "next/navigation";

export default function Home() {
  // Assim que alguém acessar localhost:3000/, o Next.js redireciona para o login
  redirect("/login");
}