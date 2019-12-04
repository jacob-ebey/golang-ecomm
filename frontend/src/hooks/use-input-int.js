import useInput from "@rooks/use-input";

export default function useInputInt(defaultValue?: string) {
  const input = useInput(defaultValue, {
    validate: newValue => {
      if (newValue === "") {
        return true;
      }

      const parsed = Number.parseFloat(newValue, 10);
      return Number.isInteger(parsed) && parsed > 0;
    }
  });

  const parsed = Number.parseFloat(input.value, 10);

  return [input, Number.isInteger(parsed) ? parsed : 0];
}
