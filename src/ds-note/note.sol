/// note.sol -- the `note' modifier, for logging calls as events

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

pragma solidity >=0.4.23;

contract DSNote {
    event LogNote(
        bytes4   indexed  sig,
        address  indexed  guy,
        bytes32  indexed  foo,
        bytes32  indexed  bar,
        uint256           wad,
        bytes             fax
    ) anonymous;

    modifier note {
        bytes32 foo;
        bytes32 bar;
        uint256 wad;

        assembly {
            foo := calldataload(4)
            bar := calldataload(36)
            wad := callvalue()
        }

        _;

        emit LogNote(msg.sig, msg.sender, foo, bar, wad, msg.data);
    }
}


contract DSNoteTest is DSTest, DSNote {
    function test_0() public { this.foo(); }
    function test_1() public { this.foo(1); }
    function test_2() public { this.foo(1, 2); }
    function test_3() public { this.foo(1, 2, 3); }
    function test_4() public { this.bar.value(uint(0x1234))(); }

    function foo() public note {}
    function foo(uint a) public note { a; }
    function foo(uint a, uint b) public note { a; b; }
    function foo(uint a, uint b, uint c) public note { a; b; c; }
    function bar() public payable note {}
}

