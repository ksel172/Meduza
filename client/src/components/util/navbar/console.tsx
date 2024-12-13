/**
 * v0 by Vercel.
 * @see https://v0.dev/t/Dp7EuNDak2F
 * Documentation: https://v0.dev/docs#integrating-generated-code-into-your-nextjs-app
 */
export default function ConsoleWidget() {
    return (
        // flex flex-col h-[100%]
      <div className="flex flex-col h-[100%] bg-[#1e1e1e] font-mono">
        <div className="flex-col h-[100%] overflow-y-auto p-1 grid gap-4">
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <span className="text-[#9cdcfe]">user@terminal</span>
              <span className="text-[#ce9178]">~</span>
              <span className="animate-blink text-[#d4d4d4]">_</span>
            </div>
            <div>
              <span className="text-[#9cdcfe]">$ </span>
              <span>echo "Hello, world!"</span>
            </div>
            <div>
              <span className="text-[#9cdcfe]">$ </span>
              <span>ls -la</span>
            </div>
            <div>
              <span className="text-[#9cdcfe]">$ </span>
              <span>cd /usr/bin</span>
            </div>
            <div>
              <span className="text-[#9cdcfe]">$ </span>
              <span>python3 --version</span>
            </div>
            <div>
              <span className="text-[#9cdcfe]">$ </span>
              <span>exit</span>
            </div>
          </div>
        </div>
        <div className="bg-[#2d2d2d] sticky bottom-0 bg-background p-4 flex items-center gap-2">
          <span className="text-[#9cdcfe]">user@terminal</span>
          <span className="text-[#ce9178]">~</span>
          <span className="animate-blink text-[#d4d4d4]">_</span>
          <input type="text" className="bg-transparent outline-none flex-1" placeholder="Enter command" />
        </div>
      </div>
    )
  }