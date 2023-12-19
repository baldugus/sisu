

function CardWrapper({children}:any) {
  return (
    <div className="shadow-md rounded-xl">
      <div className="px-2 shadow-[inset_0_-4px_0px_rgba(200,200,200,1)] py-4 items-center bg-gray-100  rounded-xl ">
        {children}
      </div>
    </div>
  )
}

export default CardWrapper