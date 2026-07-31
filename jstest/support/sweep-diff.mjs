import { readFileSync, readdirSync } from "node:fs";
import { join } from "node:path";

const [baseDir, curDir] = process.argv.slice(2);
let differences = 0;

for (const file of readdirSync(baseDir).filter((f) => f.endsWith(".json"))) {
  const before = JSON.parse(readFileSync(join(baseDir, file), "utf8"));
  const after = JSON.parse(readFileSync(join(curDir, file), "utf8"));
  for (const [element, props] of Object.entries(before)) {
    if (!(element in after)) {
      console.log(`${file}  ${element}  DISAPPEARED`);
      differences++;
      continue;
    }
    for (const [prop, value] of Object.entries(props)) {
      const now = after[element][prop];
      if (now !== value) {
        console.log(`${file}  ${element}  ${prop}\n    before: ${value}\n    after:  ${now}`);
        differences++;
      }
    }
  }
  for (const element of Object.keys(after)) {
    if (!(element in before)) {
      console.log(`${file}  ${element}  APPEARED`);
      differences++;
    }
  }
}

if (differences > 0) {
  console.error(`\n${differences} computed-style difference(s) — each needs a fix or a justification.`);
  process.exit(1);
}
console.log("no computed-style differences");
