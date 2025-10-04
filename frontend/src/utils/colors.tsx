import { theme } from "../theme/theme";

export const cores = [
  theme.palette.primary.dark,
  theme.palette.primary.light,
  theme.palette.secondary.main,
  theme.palette.secondary.light,
  theme.palette.secondary.dark,
];

/**
 * Converte uma cor hexadecimal para RGB.
 * @param {string} hex - Cor no formato "#RRGGBB".
 * @returns {[number, number, number]} Array com valores RGB.
 * @example
 * hexToRgb("#EA4700"); // [234, 71, 0]
 */
export function hexToRgb(hex: string): [number, number, number] {
  const bigint = parseInt(hex?.replace("#", ""), 16);
  const r = (bigint >> 16) & 255;
  const g = (bigint >> 8) & 255;
  const b = bigint & 255;
  return [r, g, b];
}

/**
 * Converte uma cor RGB para hexadecimal.
 * @param {{ r: number; g: number; b: number }} param0 - Objeto com valores RGB.
 * @returns {string} Cor no formato "#RRGGBB".
 * @example
 * rgbToHex({ r: 234, g: 71, b: 0 }); // "#EA4700"
 */
export function rgbToHex({
  r,
  g,
  b,
}: {
  r: number;
  g: number;
  b: number;
}): string {
  return (
    "#" +
    ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1).toUpperCase()
  );
}

/**
 *  Interpola entre duas cores hexadecimais.
 * @param {{ c1: string; c2: string; t: number }} param0 - Objeto com as cores e o fator de interpolação.
 * @returns {string} Cor interpolada no formato "#RRGGBB".
 * @example
 * interpolateColor({ c1: "#EA4700", c2: "#FFE46A", t: 0.5 }); // "#F28B35"
 */
export function interpolateColor({
  c1,
  c2,
  t,
}: {
  c1: string;
  c2: string;
  t: number;
}): string {
  const [r1, g1, b1] = hexToRgb(c1);
  const [r2, g2, b2] = hexToRgb(c2);
  const r = Math.round(r1 + (r2 - r1) * t);
  const g = Math.round(g1 + (g2 - g1) * t);
  const b = Math.round(b1 + (b2 - b1) * t);
  return rgbToHex({ r, g, b });
}

/**
 * Mapeia um valor para uma cor em um gradiente definido por um array de cores.
 * @param {{ cores: string[]; valor: number | string; min?: number; max?: number }} param0 - Objeto com as cores, valor e limites.
 * @returns {string} Cor mapeada no formato "#RRGGBB".
 * @example
 * mapearCor({ cores: ["#EA4700", "#FFE46A", "#8CC0B7"], valor: 5, min: 1, max: 10 }); // "#FFB24D"
 * mapearCor({ cores: ["#EA4700", "#FFE46A", "#8CC0B7"], valor: "∞" }); // "#0B1926"
 */
export function mapearCor({
  valor,
  min = 1,
  max = 55,
}: {
  valor: number | string;
  min?: number;
  max?: number;
}): string {
  if (typeof valor === "number") {
    const t = (valor - min) / (max - min); // Normaliza entre 0 e 1
    const segment = t * (cores.length - 1);
    const index = Math.floor(segment);
    const localT = segment - index;
    return interpolateColor({
      c1: cores[index],
      c2: cores[index + 1],
      t: localT,
    });
  } else {
    return "#111";
  }
}
