/**
 * v0 by Vercel.
 * @see https://v0.dev/t/Dp7EuNDak2F
 * Documentation: https://v0.dev/docs#integrating-generated-code-into-your-nextjs-app
 */
export default function ConsoleWidget() {
  return (
    // flex flex-col h-[100%]
    <div className="flex flex-col h-[100%] w-[100%] bg-[#0F172A] rounded-lg">
      <div className="flex-col h-[100%] overflow-y-auto p-1 grid gap-4">
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <span className="text-[#9cdcfe]">voldemort@theagency:</span>
            <span className="text-[#ce9178]">~$</span>
            <span>cmd /c dir</span>
            {/* <span className="animate-blink text-[#d4d4d4]">_</span> */}
          </div>
          <div>
            {/* <span className="text-[#9cdcfe]">$ </span> */}
            <span>formatted response</span>
          </div>
        </div>
      </div>
      <div className="sticky bottom-0 bg-[#0F172A] p-4 flex items-center gap-2 rounded-lg">
        <span className="text-[#9cdcfe]">voldemort@theagency:</span>
        <span className="text-[#ce9178]">~$</span>
        <input
          type="text"
          className="bg-[#0F172A] outline-none flex-1 bottom-3 rounded-lg"
          placeholder="Enter command"
        />
      </div>
    </div>
  );
}
