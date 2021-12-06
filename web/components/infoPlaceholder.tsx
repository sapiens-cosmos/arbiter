export default function InfoPlaceholder({ className }: { className: string }) {
  return (
    <div
      className={`${className} relative overflow-hidden rounded-sm my-0.25 bg-gray bg-opacity-60`}
    >
      <div className="placeholder-animation" />
    </div>
  );
}
