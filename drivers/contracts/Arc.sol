pragma solidity 0.4.23;

contract Arc {
    mapping (bytes32=>bytes) public swaps;

    function submit(bytes32 _orderID, bytes _swapDetails) public {
        swaps[_orderID] = _swapDetails;
    }
}